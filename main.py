from fastapi import FastAPI
from pydantic import BaseModel
from dotenv import load_dotenv
import os
import requests

load_dotenv()
API_KEY = os.getenv("GEMINI_API_KEY")
GEMINI_ENDPOINT = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent"

app = FastAPI()

class SuggestRequest(BaseModel):
    keywords: list[str]
    tone: str

def parse_suggestion(text: str) -> dict:
    lines = text.split('\n')
    data = {
        "title": "",
        "audience": "",
        "headlines": []
    }
    for i, line in enumerate(lines):
        if "Blog Post Idea:" in line:
            data["title"] = line.replace("Blog Post Idea:", "").strip()
        elif "Target Audience:" in line:
            data["audience"] = line.replace("Target Audience:", "").strip()
        elif line.startswith("* "):
            data["headlines"].append(line[2:].strip())
    return data

@app.post("/generate")
def generate_blog_suggestion(body: SuggestRequest):
    prompt_text = f"Write a {body.tone} blog post idea using these keywords: {', '.join(body.keywords)}"
    
    headers = {
        "Content-Type": "application/json",
        "X-goog-api-key": API_KEY
    }

    payload = {
        "contents": [
            {
                "parts": [
                    {"text": prompt_text}
                ]
            }
        ]
    }

    response = requests.post(GEMINI_ENDPOINT, headers=headers, json=payload)
    data = response.json()

    try:
        text = data["candidates"][0]["content"]["parts"][0]["text"]
    except (KeyError, IndexError):
        text = "Failed to generate content. Please try again."
        return {"error": text}

    return parse_suggestion(text)
