import os
from dotenv import load_dotenv
load_dotenv()
class Config:
    KAFKA_BOOTSTRAP_SERVERS = "localhost:9092"
    REC_BEATS_TOPIC = "rec_beats_topic2"
    REFILL_TOPIC = "rec_refill_requests"
    REFILL_COOLDOWN = 300
    REFILL_THRESHOLD = 5  # Пороговое значение для дозапроса
    REFILL_COUNT = 9     
    BATCH_SIZE = 9       # Максимальное кол-во рекомендаций за раз
    MIN_GENRES = 1       
    MAX_GENRES = 3
    WT_SECRET_KEY = os.getenv("SECRET_KEY")  # Берём из переменной SECRET_KEY в .env
    JWT_TOKEN_LOCATION = ["headers"]          # Искать токен в заголовках Authorization
    JWT_ACCESS_TOKEN_EXPIRES = 3600      