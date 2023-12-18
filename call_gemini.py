# -*- coding: utf-8 -*-

import google.generativeai as genai
import sys

sys.stdout.reconfigure(encoding='utf-8')

with open('key.txt') as f:
    API_KEY = f.read().strip()

with open('req.txt', encoding="utf-8") as f:
    REQ = f.read().strip()

genai.configure(api_key=API_KEY)

# Set up the model
generation_config = {
    "temperature": 0,
    "top_p": 1,
    "top_k": 1,
    "max_output_tokens": 2048,
}


model = genai.GenerativeModel(model_name="gemini-pro",
                              generation_config=generation_config)

prompt_parts = [
    REQ
]

response = model.generate_content(prompt_parts)
print(response.text)
