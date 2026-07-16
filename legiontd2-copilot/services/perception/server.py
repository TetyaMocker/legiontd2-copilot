"""Perception Service — единственный компонент, работающий с изображениями напрямую.

Phase 1.3: OCR Mythium/HP/таймер/номер волны через mss + EasyOCR.
Phase 2.3: CV-детекция юнитов.
"""

import logging
import time
from concurrent import futures

import grpc
import mss
import numpy as np

import perception_pb2
import perception_pb2_grpc

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
log = logging.getLogger("perception")


class PerceptionServicer(perception_pb2_grpc.PerceptionServiceServicer):
    def __init__(self):
        self.sct = mss.mss()
        self._last_health = time.time()

    def _grab_screen(self):
        mon = self.sct.monitors[1]
        return np.array(self.sct.grab(mon))

    def ReadEconomy(self, request, context):
        # Phase 1.3 stub: распознавание через OCR
        log.info("ReadEconomy called")
        return perception_pb2.EconomyState(
            mythium=0,
            income=0,
            wave_number=0,
            wave_timer_seconds=0,
            king_hp_percent=0,
            ally_king_hp_percent=0,
            confidence=0.0,
        )

    def ReadBattlefield(self, request, context):
        log.info("ReadBattlefield called — not implemented until Phase 2.3")
        return perception_pb2.BattlefieldState(confidence=0.0)

    def HealthCheck(self, request, context):
        return perception_pb2.HealthCheckResponse(
            game_window_detected=True,
            capture_healthy=True,
            message="Perception Service running (Phase 1.3 skeleton)",
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
