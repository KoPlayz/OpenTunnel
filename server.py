# fileserver.py
from http.server import SimpleHTTPRequestHandler, HTTPServer
import json, os

class MyHandler(SimpleHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/list":
            files = [f for f in os.listdir(".") if f.endswith(".ote")]
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(json.dumps(files).encode())
        else:
            super().do_GET()

if __name__ == "__main__":
    server = HTTPServer(("0.0.0.0", 8080), MyHandler)
    print("Serving on port 8080")
    server.serve_forever()
