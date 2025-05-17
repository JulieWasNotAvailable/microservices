import uuid
from flask import Blueprint, jsonify, request
from services.s3_services import check_file_in_s3
from services.kafka_service import send_kafka_message
from services.audio_service import analyze_audio
from werkzeug.utils import secure_filename
import os

bp = Blueprint('api', __name__)

BUCKET_NAME = os.getenv('BUCKET_NAME', 'music-beats-bucket')
MP3_FOLDER = os.getenv('MP3_FOLDER', 'mp3/')
KAFKA_TRACK_TOPIC = os.getenv('KAFKA_TRACK_TOPIC', 'track_for_mfcc')

@bp.route('/api/process_tracks', methods=['POST'])
def process_tracks():
    """Обработка одного или нескольких треков"""
    if not request.json:
        return jsonify({"error": "No data provided"}), 400

    filenames = request.json if isinstance(request.json, list) else [request.json]

    processed = []
    errors = []

    for filename in filenames:
        try:
            # Проверяем существование файла в S3
            if not check_file_in_s3(filename):
                errors.append({"filename": filename, "error": "File not found in S3"})
                continue

            # Генерация beat_id
            beat_id = str(uuid.uuid4())

            message = {"filename": filename, "beat_id": beat_id}

            if send_kafka_message(KAFKA_TRACK_TOPIC, message):
                processed.append({"filename": filename, "beat_id": beat_id, "status": "queued"})
            else:
                errors.append({"filename": filename, "error": "Failed to send to Kafka"})

        except Exception as e:
            errors.append({"filename": filename, "error": str(e)})

    return jsonify({"processed": processed, "errors": errors, "message": f"Processed {len(processed)} tracks, {len(errors)} errors"}), 200
