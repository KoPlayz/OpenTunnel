from http.server import BaseHTTPRequestHandler, HTTPServer
import json, os, base64, urllib.request, urllib.parse
from datetime import datetime

PORT = 8080
API_KEY = "ADadm7QY50k2rySCj0Nang69hhZne8SR"  # Replace with your desired API key

LOG_DIR = "Logs"
os.makedirs(LOG_DIR, exist_ok=True)  # Ensure the Logs directory exists

class MyHandler(BaseHTTPRequestHandler):
    def log_request(self, message):
        """Log the request to a file."""
        log_file = os.path.join(LOG_DIR, f"{datetime.now().date().isoformat()}.log")
        with open(log_file, "a") as log:
            log.write(f"{datetime.now().isoformat()} - {self.client_address[0]} - {message}\n")

    def authenticate(self):
        """Check if the request contains a valid API key."""
        api_key = self.headers.get("X-API-Key")  # Read the API key from the request headers
        if api_key != API_KEY:
            self.log_request("Unauthorized access attempt")
            self.send_response(401)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(json.dumps({"error": "Unauthorized"}).encode())
            return False
        return True

    def do_GET(self):
        self.log_request(f"GET {self.path}")
        if not self.authenticate():
            return  # Stop processing if authentication fails

        if self.path == "/list":
            # Encode filenames to ensure they are URL-safe
            files = [urllib.parse.quote(f) for f in os.listdir(".") if f.endswith(".ote")]
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(json.dumps(files).encode())
        else:
            path = urllib.parse.unquote(self.path.lstrip("/"))  # Decode the requested filename
            if os.path.exists(path):
                self.send_response(200)
                self.send_header("Content-Type", "application/octet-stream")
                self.end_headers()
                with open(path, "rb") as f:
                    self.wfile.write(f.read())
            else:
                self.send_error(404, "File not found")

    def do_POST(self):
        self.log_request(f"POST {self.path}")
        if not self.authenticate():
            return  # Stop processing if authentication fails

        if self.path == "/add":
            length = int(self.headers['Content-Length'])
            body = self.rfile.read(length)
            data = json.loads(body.decode())

            url = data.get("url")
            if not url:
                self.send_error(400, "Missing url")
                return

            self.log_request(f"URL requested: {url}")  # Log the requested URL

            try:
                # Open URL and follow redirects to actual file
                with urllib.request.urlopen(url) as response:
                    # Get the final URL after redirects
                    final_url = response.geturl()

                    # Try to get the filename from the Content-Disposition header
                    content_disposition = response.headers.get("Content-Disposition")
                    if content_disposition and "filename=" in content_disposition:
                        filename = content_disposition.split("filename=")[-1].strip('"')
                    else:
                        # Fallback to the final URL path
                        filename = os.path.basename(urllib.parse.urlparse(final_url).path) or "file"

                    # Replace %20 with a space in the filename (URL decoding)
                    filename = urllib.parse.unquote(filename)

                    # Save the file
                    with open(filename, "wb") as f:
                        f.write(response.read())

                # Encode to base64
                with open(filename, "rb") as f:
                    encoded = base64.b64encode(f.read()).decode()

                # Save as .ote
                ote_name = filename + ".ote"
                with open(ote_name, "w") as f:
                    f.write(encoded)

                os.remove(filename)  # clean original

                self.send_response(200)
                self.send_header("Content-Type", "application/json")
                self.end_headers()
                self.wfile.write(json.dumps({"status": "ok", "file": ote_name}).encode())

            except Exception as e:
                self.log_request(f"Error processing POST /add: {str(e)}")
                self.send_error(500, str(e))

if __name__ == "__main__":
    server = HTTPServer(("0.0.0.0", PORT), MyHandler)
    print(f"Serving on port {PORT}")
    server.serve_forever()