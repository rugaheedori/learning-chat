import WebSocket from "ws";
import { socketLogger } from "../logs/winston";

class Room {
  private clients: Set<WebSocket>;

  constructor() {
    this.clients = new Set();
  }

  join(client: WebSocket) {
    socketLogger.info("new client");
    this.clients.add(client);
  }

  leave(client: WebSocket) {
    socketLogger.info("removed client");
    this.clients.delete(client);
  }

  forwardMessage(message: {name: string, msg: string}) {
    socketLogger.info("send message all clients");
    for (const client of this.clients) {
      client.send(JSON.stringify(message));
    }
  }
}

export { Room };
