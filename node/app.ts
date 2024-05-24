import http from "http";
import Websocket from "ws";
import { socketLogger } from "./logs/winston";
import { Room } from "./types/room";
import { Request } from "express";

const server = http.createServer();
const wss = new Websocket.Server({ server });

const room = new Room();

wss.on("connection", (ws: Websocket, req: Request) => {
  // Cookie에서 user 정보 가지고 오기
  const cookie = req.headers.cookie;
  const [_, user] = (cookie as string).split("=");

  room.join(ws);

  ws.on("message", (msg: string) => {
    // message가 들어오면 해당 메세제를 다른 client에도 브로드 캐스팅
    const jsonMsg = JSON.parse(msg);

    jsonMsg.name = user;

    room.forwardMessage(jsonMsg);
  });

  ws.on("close", () => {
    // client 접속 끊긴 경우, client 제거
    room.leave(ws);
  });
});

const PORT = 8000;

server.listen(PORT, () => {
  socketLogger.info(`Server started on PORT: ${PORT}`);
});
