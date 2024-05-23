import { createLogger, format, transports } from "winston";
import winstonDaily from "winston-daily-rotate-file";

const logDri = "log";

const logFormat = format.printf((info) => {
  return `${info.timestmap} ${info.level} ${info.message}`;
});

const socketLogger = createLogger({
  format: format.combine(
    format.timestamp({
      format: "YYYY-MM-DD HH:mm:ss",
    }),
    logFormat
  ),
  transports: [
    new winstonDaily({
      filename: "socket-info.log",
      datePattern: "YYYY-MM-DD",
      dirname: logDri + "/socket",
      level: "info",
    }),
    new winstonDaily({
      filename: "socket-err.log",
      datePattern: "YYYY-MM-DD",
      dirname: logDri + "/socket",
      level: "error",
    }),
  ],
});

socketLogger.add(
  new transports.Console({
    format: format.combine(format.colorize(), format.simple()),
  })
);

export {socketLogger}