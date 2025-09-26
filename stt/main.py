import sys
import time
from pathlib import Path
from audio_transcriber import AudioTranscriber 

FILE_ARG="-f" # path to input audio file
DIR_ARG="-d"  # path to input audio directory
WATCH_ARG="-w" # watch directory for new files
WATCH_PROCESSED_ARG="-wo" # directory to move processed files to

# # Usage example
# def transcribe_audio_file(audio_path: str, model_size="base", language=None):
#     """Simple function to transcribe an audio file"""
    
#     transcriber = AudioTranscriber(model_size=model_size)
#     result = transcriber.transcribe_file(audio_path, language=language)
    
#     # Save transcription
#     audio_name = Path(audio_path).stem
#     output_path = f"{audio_name}_transcript.txt"
#     transcriber.save_transcription(result, output_path)
    
#     return result

# Usage example
def transcribe_audio_file(ts: AudioTranscriber, audio_path: str, output_path:str, language=None):
    """Simple function to transcribe an audio file"""
    
    result = ts.transcribe_file(audio_path, language=language)
    ts.save_transcription(result, output_path)
    
    return result


def get_arg_value(arg_name: str, usage_msg: str) -> str:
    """Check command line arguments for a specific argument and return its value"""
    if arg_name in sys.argv:
        file_index = sys.argv.index(arg_name) + 1
        if file_index < len(sys.argv):
            return sys.argv[file_index]
        else:
            print(f"Error: No value provided after {arg_name}")
            sys.exit(1)

def watch_directory(path: str, processed_out: str = None):
    """Watch a directory for new audio files and transcribe them as they appear"""
    print(f"Watching directory: {path}")
    processed_files = set()
    iterations_since_last_processed = 0
    transcriber = AudioTranscriber(model_size="small")
    
    while True:
        iterations_since_last_processed += 1
        current_files = set(Path(path).glob("*"))
        new_files = current_files - processed_files

        for file in new_files:
            if file.suffix.lower() in ['.wav', '.mp3', '.m4a', '.flac', '.aac', '.ogg']:
                print(f"\nNew file detected: {file.name}")
                output_file_path = str(file.with_suffix(''))
                if processed_out:
                    output_file_path = f"{processed_out}/{Path(file).stem}"
                output_path = f"{output_file_path}_transcript.txt"
                transcribe_audio_file(transcriber, str(file), output_path, language="en")
                processed_files.add(file)
                iterations_since_last_processed = 0

                if processed_out:
                    processed_dir = Path(processed_out)
                    processed_dir.mkdir(parents=True, exist_ok=True)
                    dest_path = processed_dir / file.name
                    file.rename(dest_path)
                    print(f"Moved processed file to: {dest_path}")

        # Sleep briefly to avoid busy waiting
        time.sleep(1)

        if iterations_since_last_processed >= 60: # No new files for 60 iterations (>= 1 minute)
            time.sleep(30)  # Sleep longer if no new files were processed

# Example usage
if __name__ == "__main__":
    file_arg = get_arg_value(FILE_ARG, f"Usage: python {sys.argv[0]} {FILE_ARG} <path_to_audio_file>")
    dir_arg = get_arg_value(DIR_ARG, f"Usage: python {sys.argv[0]} {DIR_ARG} <path_to_audio_directory>")
    watch_arg = get_arg_value(WATCH_ARG, f"Usage: python {sys.argv[0]} {WATCH_ARG} <path_to_watch_directory>")
    processed_out_arg = get_arg_value(WATCH_PROCESSED_ARG, f"Usage: python {sys.argv[0]} {WATCH_PROCESSED_ARG} <path_to_move_processed_files>")

    if not file_arg and not dir_arg and not watch_arg:
        print(f"Usage: python {sys.argv[0]} {FILE_ARG} <path_to_audio_file> OR {DIR_ARG} <path_to_audio_directory> OR {WATCH_ARG} <path_to_watch_directory>")
        sys.exit(1)

    if watch_arg:
        watch_directory(watch_arg, processed_out_arg)
        sys.exit(0)

    files = []
    if file_arg:
        files.append(file_arg)
    if dir_arg:
        dir_path = Path(dir_arg)
        if dir_path.is_dir():
            files.extend([str(f) for f in dir_path.glob("*")])
            print(f"Found {len(files)} files in directory {dir_arg}")
        else:
            print(f"Error: {dir_arg} is not a valid directory.")
            sys.exit(1)

    transcriber = AudioTranscriber(model_size="small")

    for file in files:
        if file.endswith('.txt'):
            print(f"Skipping text file: {file}")
            continue
        print(f"\nTranscribing file: {file}")
        fullpath_without_ext = str(Path(file).with_suffix(''))
        output_path = f"{fullpath_without_ext}_transcript.txt"
        result = transcribe_audio_file(transcriber, file, output_path, language="en")

    print(f"\nTranscription preview:")
    print(result['text'][:200] + "..." if len(result['text']) > 200 else result['text'])
