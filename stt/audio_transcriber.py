# Code originally written by Pavlin
# Retrieved from: https://www.pavlinbg.com/posts/python-speech-to-text-guide
# Retrieved on: 2026-09-25

import whisper
import os
from pathlib import Path
import time

class AudioTranscriber:
    def __init__(self, model_size="base"):
        """Initialize transcriber with specified Whisper model"""
        print(f"Loading Whisper {model_size} model...")
        self.model = whisper.load_model(model_size)
        print("Model loaded successfully!")
    
    def transcribe_file(self, audio_path, language=None):
        """
        Transcribe a single audio file
        
        Args:
            audio_path: Path to audio file
            language: Language code ('en', 'es', 'fr', etc.) or None for auto-detect
        """
        if not os.path.exists(audio_path):
            raise FileNotFoundError(f"Audio file not found: {audio_path}")
        
        print(f"Transcribing: {Path(audio_path).name}")
        
        start_time = time.time()
        
        # Transcribe audio
        options = {"language": language} if language else {}
        result = self.model.transcribe(audio_path, **options)
        
        processing_time = time.time() - start_time
        
        print(f"✓ Completed in {processing_time:.1f} seconds")
        print(f"✓ Detected language: {result['language']}")
        
        return {
            'text': result['text'].strip(),
            'language': result['language'],
            'segments': result.get('segments', []),
            'processing_time': processing_time
        }
    
    def save_transcription(self, result, output_path):
        """Save transcription to text file"""
        with open(output_path, 'w', encoding='utf-8') as f:
            f.write("=== Transcription Results ===\n")
            f.write(f"Language: {result['language']}\n")
            f.write(f"Processing Time: {result['processing_time']:.1f} seconds\n")
            f.write("=" * 40 + "\n\n")
            f.write(result['text'])
        
        print(f"✓ Transcription saved to: {output_path}")
