import express from "express";
import http from "http";
import cors from "cors";
import Websocket from "ws";

const app = express();

app.use(
  cors({
    origin: "*",
  })
);
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
