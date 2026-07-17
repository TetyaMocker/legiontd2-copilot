"""Perception Service — screen capture + OCR via EasyOCR, POSTs to Go orchestrator."""

import json
import logging
import os
import re
import time
import urllib.request

import cv2
import mss
import numpy as np

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
log = logging.getLogger("perception")

ORCHESTRATOR_URL = os.getenv("ORCHESTRATOR_URL", "http://localhost:8080/api/ingest")
OCR_INTERVAL = float(os.getenv("OCR_INTERVAL", "2.0"))

# Configuration regions matching config/regions.json
REGIONS = {
    "mythium":  (520, 50, 120, 35, 0, 99999),
    "income":   (515, 108, 80, 28, 0, 9999),
    "wave":     (850, 90, 160, 35, 1, 30),
    "timer":    (890, 58, 100, 30, 0, 600),
    "king_hp":  (1410, 85, 70, 35, 0, 100),
}


def init_reader():
    import easyocr
    log.info("initializing EasyOCR reader...")
    reader = easyocr.Reader(["en"], gpu=False, verbose=False)
    log.info("EasyOCR reader ready")
    return reader


def ocr_number(reader, img, x, y, w, h, vmin=0, vmax=99999):
    crop = img[y: y + h, x: x + w]
    if crop.size == 0:
        return None
    gray = np.dot(crop[..., :3], [0.299, 0.587, 0.114]).astype(np.uint8)
    _, thresh = cv2.threshold(gray, 0, 255, cv2.THRESH_BINARY + cv2.THRESH_OTSU)
    results = reader.readtext(thresh, detail=0, paragraph=False)
    for text in results:
        digits = re.sub(r"[^0-9]", "", text)
        if digits:
            val = int(digits)
            if vmin <= val <= vmax:
                return val
    return None


class Tracker:
    def __init__(self):
        self.cache = {}
        self.last_scan = 0

    def should_scan(self, now, interval):
        if now - self.last_scan >= interval:
            return True
        return False

    def update(self, key, value):
        self.cache[key] = {"value": value, "time": time.time()}

    def get(self, key):
        entry = self.cache.get(key)
        if entry and time.time() - entry["time"] < 10:
            return entry["value"]
        return None


def capture_once(reader):
    sct = mss.MSS()
    img = np.array(sct.grab(sct.monitors[1]))

    result = {"confidence": 0.0}
    found = 0

    for name, (x, y, w, h, vmin, vmax) in REGIONS.items():
        val = ocr_number(reader, img, x, y, w, h, vmin, vmax)
        if val is not None:
            result[name] = val
            found += 1
        else:
            result[name] = 0

    if found >= 3:
        result["confidence"] = min(0.3 + found * 0.15, 1.0)

    return result


def main():
    reader = init_reader()
    tracker = Tracker()

    log.info("starting capture loop, posting to %s", ORCHESTRATOR_URL)
    while True:
        now = time.time()
        if tracker.should_scan(now, OCR_INTERVAL):
            state = capture_once(reader)
            tracker.last_scan = now

            for key, val in state.items():
                if key != "confidence":
                    tracker.update(key, val)

            log.info("state: m=%d inc=%d hp=%d wave=%d timer=%d conf=%.2f",
                     state.get("mythium", 0), state.get("income", 0),
                     state.get("king_hp", 0), state.get("wave", 0),
                     state.get("timer", 0), state.get("confidence", 0))

            try:
                data = json.dumps(state).encode("utf-8")
                req = urllib.request.Request(
                    ORCHESTRATOR_URL,
                    data=data,
                    headers={"Content-Type": "application/json"},
                )
                resp = urllib.request.urlopen(req, timeout=3)
                resp.read()
            except Exception as e:
                log.warning("failed to POST to orchestrator: %s", e)
        else:
            time.sleep(0.1)


if __name__ == "__main__":
    main()
