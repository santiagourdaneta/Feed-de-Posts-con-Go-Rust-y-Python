from fastapi import FastAPI
from textblob import TextBlob
from pydantic import BaseModel

# 1. Le decimos a Python que cree una nueva puerta mágica (API)
app = FastAPI()

# 2. Creamos una receta de lo que vamos a recibir
class PostText(BaseModel):
    content: str
    post_id: int

# 3. Le decimos al mago que se pare en una puerta para escuchar
@app.post("/analyze/sentiment")
async def analyze_sentiment(post: PostText):
    # 4. Cuando alguien le dé un texto, el mago TextBlob lo analiza
    analysis = TextBlob(post.content)
    sentiment_score = analysis.sentiment.polarity
    
    # 5. El mago decide el sentimiento y le pone una etiqueta
    if sentiment_score > 0:
        sentiment_label = "positivo"
    elif sentiment_score < 0:
        sentiment_label = "negativo"
    else:
        sentiment_label = "neutro"

    # 6. Muestra el resultado de la magia en la consola
    print(f"Post ID: {post.post_id} - Sentimiento: {sentiment_label}")

    # 7. Finalmente, le devuelve el resultado al que se lo pidió (a Go)
    return {"post_id": post.post_id, "sentiment": sentiment_label}