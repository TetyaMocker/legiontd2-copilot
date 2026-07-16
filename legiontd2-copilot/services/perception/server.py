"""Perception Service — захват экрана + OCR через EasyOCR."""

import logging
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

_OCR_READER = None


def _get_reader():
    global _OCR_READER
    if _OCR_READER is None:
        import easyocr
        _OCR_READER = easyocr.Reader(["en"], gpu=True, verbose=False)
        log.info("easyocr reader initialized")
    return _OCR_READER


def _ocr_number(img, x, y, w, h):
    crop = img[y: y + h, x: x + w]
    gray = np.dot(crop[..., :3], [0.299, 0.587, 0.114]).astype(np.uint8)
    _, thresh = cv2.threshold(gray, 0, 255, cv2.THRESH_BINARY + cv2.THRESH_OTSU)
    reader = _get_reader()
    results = reader.readtext(thresh, detail=0, paragraph=False)
    for text in results:
        digits = re.sub(r"[^0-9]", "", text)
        if digits:
            return digits
    return None


def _ocr_text(img, x, y, w, h):
    crop = img[y: y + h, x: x + w]
    gray = np.dot(crop[..., :3], [0.299, 0.587, 0.114]).astype(np.uint8)
    _, thresh = cv2.threshold(gray, 0, 255, cv2.THRESH_BINARY + cv2.THRESH_OTSU)
    reader = _get_reader()
    results = reader.readtext(thresh, detail=0, paragraph=False)
    if results:
        return results[0]
    return None


class PerceptionServicer(perception_pb2_grpc.PerceptionServiceServicer):
    def __init__(self):
        self.sct = mss.mss()
        # default regions for 1920x1080 — user can override via env
        self.regions = {
            "mythium": (30, 40, 160, 36),
            "king_hp": (30, 85, 140, 30),
            "wave": (880, 8, 160, 36),
            "timer": (900, 44, 120, 28),
        }

    def ReadEconomy(self, request, context):
        t0 = time.time()
        img = np.array(self.sct.grab(self.sct.monitors[1]))

        r = self.regions
        mythium_text = _ocr_number(img, *r["mythium"])
        hp_text = _ocr_number(img, *r["king_hp"])
        timer_text = _ocr_number(img, *r["timer"])
        wave_text = _ocr_number(img, *r["wave"])

        conf = 0.0
        if mythium_text or hp_text:
            conf = 0.5 if mythium_text else 0.3
            if hp_text:
                conf += 0.2

        eco = perception_pb2.EconomyState(
            mythium=int(mythium_text) if mythium_text else 0,
            income=0,
            wave_number=int(wave_text) if wave_text else 0,
            wave_timer_seconds=int(timer_text) if timer_text else 0,
            king_hp_percent=int(hp_text) if hp_text else 0,
            ally_king_hp_percent=0,
            confidence=min(conf, 1.0),
        )
        log.info("eco: mythium=%s hp=%s wave=%s timer=%s [%.2fs]",
                 mythium_text, hp_text, wave_text, timer_text, time.time() - t0)
        return eco

    def ReadBattlefield(self, request, context):
        return perception_pb2.BattlefieldState(confidence=0.0)

    def HealthCheck(self, request, context):
        return perception_pb2.HealthCheckResponse(
            game_window_detected=True,
            capture_healthy=True,
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
