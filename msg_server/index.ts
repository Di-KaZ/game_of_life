
import { uniqueNamesGenerator, adjectives, colors, animals } from 'unique-names-generator';


const sockets: ServerWebSocket[] = []


interface Payload {
	name: string
	content: string
}

interface SocketIdentity {
	name: string
}

Bun.serve<SocketIdentity>({
	hostname: "0.0.0.0",
	port: 8080,
	fetch(req, server) {
		const name = uniqueNamesGenerator({ dictionaries: [adjectives, colors, animals] });
		if (server.upgrade(req, { data: { name } })) return;
		return new Response("Upgrade failed :(", { status: 500 });
	},
	websocket: {
		open(ws) {
			console.info("New connection  ", ws.data.name);
			sockets.push(ws)
		},
		message(ws, message: string) {
			for (const s of sockets) {
				s.send(JSON.stringify({ name: ws.data.name, content: message } as Payload))
			}
		}
	},
});
