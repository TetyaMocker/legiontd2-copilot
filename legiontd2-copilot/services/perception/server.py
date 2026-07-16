"""Perception Service — screen capture + OCR."""

import logging
import os
import re
import time
from concurrent import futures

import cv2
import grpc
import mss
import numpy as np

import perception_pb2
import perception_pb2_grpc

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
log = logging.getLogger("perception")

# In-match HUD zones (1920x1080 default), adjusted by calibration
ZONES = {
    "mythium":  (520, 50, 120, 35),
    "income":   (515, 108, 80, 28),
    "wave":     (850, 90, 160, 35),
    "timer":    (890, 58, 100, 30),
    "king_hp":  (1410, 85, 70, 35),
}


def _ocr(reader, img, zone, valid_range=None):
    x, y, w, h = zone
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
            if valid_range and not (valid_range[0] <= val <= valid_range[1]):
                continue
            return digits
    return None


def _ocr_float(reader, img, zone, valid_range=(0, 600)):
    x, y, w, h = zone
    crop = img[y: y + h, x: x + w]
    if crop.size == 0:
        return None
    gray = np.dot(crop[..., :3], [0.299, 0.587, 0.114]).astype(np.uint8)
    _, thresh = cv2.threshold(gray, 0, 255, cv2.THRESH_BINARY + cv2.THRESH_OTSU)
    results = reader.readtext(thresh, detail=0, paragraph=False)
    for text in results:
        clean = re.sub(r"[^.0-9]", "", text)
        clean = re.sub(r"\.\.+", ".", clean)
        m = re.search(r"(\d+\.?\d*)", clean)
        if m:
            raw = m.group(1)
            val = int(float(raw))
            if valid_range and not (valid_range[0] <= val <= valid_range[1]):
                continue
            return str(val)
    return None


class PerceptionServicer(perception_pb2_grpc.PerceptionServiceServicer):
    def __init__(self):
        self.sct = mss.MSS()
        log.info("initializing EasyOCR reader...")
        import easyocr
        self.reader = easyocr.Reader(["en"], gpu=False, verbose=False)
        log.info("EasyOCR reader ready")

    def _grab(self):
        return np.array(self.sct.grab(self.sct.monitors[1]))

    def ReadEconomy(self, request, context):
        t0 = time.time()
        try:
            img = self._grab()
        except Exception as e:
            log.error("screen grab failed: %s", e)
            return perception_pb2.EconomyState(confidence=0.0)

        mythium_text = _ocr(self.reader, img, ZONES["mythium"], (0, 99999))
        income_text = _ocr(self.reader, img, ZONES["income"], (0, 9999))
        hp_text = _ocr(self.reader, img, ZONES["king_hp"], (0, 100))
        wave_text = _ocr(self.reader, img, ZONES["wave"], (1, 30))
        timer_text = _ocr_float(self.reader, img, ZONES["timer"], (0, 600))

        conf = 0.0
        if mythium_text:
            conf += 0.50
        else:
            # try to detect if game is still running from other signals
            pass
        if hp_text:
            conf += 0.25
        if wave_text:
            conf += 0.15
        if timer_text:
            conf += 0.10

        eco = perception_pb2.EconomyState(
            mythium=int(mythium_text) if mythium_text else 0,
            income=int(income_text) if income_text else 0,
            wave_number=int(wave_text) if wave_text else 0,
            wave_timer_seconds=int(timer_text) if timer_text else 0,
            king_hp_percent=int(hp_text) if hp_text else 0,
            ally_king_hp_percent=0,
            confidence=min(conf, 1.0),
        )
        log.info("eco: m=%s inc=%s hp=%s wave=%s timer=%s [%.2fs]",
                 mythium_text, income_text, hp_text, wave_text, timer_text, time.time() - t0)
        return eco

    def ReadBattlefield(self, request, context):
        return perception_pb2.BattlefieldState(confidence=0.0)

    def HealthCheck(self, request, context):
        return perception_pb2.HealthCheckResponse(
            game_window_detected=True, capture_healthy=True,
            message="Perception Service ready",
        )


def serve(port: int = 50051) -> None:
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    perception_pb2_grpc.add_PerceptionServiceServicer_to_server(PerceptionServicer(), server)
    server.add_insecure_port(f"[::]:{port}")
    log.info("Perception Service ready on :%d", port)
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
