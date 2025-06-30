import signal, sys; signal.signal(signal.SIGTERM, lambda signum, frame: sys.exit(0));
from http.server import SimpleHTTPRequestHandler, HTTPServer
class DummyHandler(SimpleHTTPRequestHandler):
	def do_GET(self):
		response_content = '{"request":{"mbean":"org.apache.activemq.artemis:broker=\\"0.0.0.0\\"","attribute":"Status","type":"read"},"value":"{\\"configuration\\":{\\"properties\\":{\\"broker.properties\\":{\\"alder32\\":\\"1\\",\\"fileAlder32\\":\\"1285359230\\"}}},\\"server\\":{\\"jaas\\":{\\"properties\\":{}},\\"state\\":\\"STARTED\\",\\"version\\":\\"2.42.0\\",\\"nodeId\\":\\"27d63164-aa5b-11f0-a502-cabd4cd943ed\\",\\"identity\\":null,\\"uptime\\":\\"1 minute\\"}}","status":200}'.encode('utf-8');
		self.send_response(200);
		self.send_header('Content-type', 'application/json');
		self.send_header('Content-Length', str(len(response_content)));
		self.end_headers();
		self.wfile.write(response_content);
HTTPServer(('', 8161), DummyHandler).serve_forever()

