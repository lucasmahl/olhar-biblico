//! ARQUIVO PARA USAR EM PRODUÇÃO
import { handler } from "./build/handler.js";
import express from "express";
import fs from "fs";
import https from "https";
import dotenv from "dotenv";

const app = express();

// Para redirecionar caso tente acessar www.olharbiblico.com
app.use(
  (
    /** @type {{ headers: { host: string; }; url: string; }} */ req,
    /** @type {{ redirect: (arg0: number, arg1: string) => void; }} */ res,
    /** @type {() => void} */ next
  ) => {
    if (req.headers.host === "www.olharbiblico.com") {
      res.redirect(301, "https://olharbiblico.com" + req.url);
    } else {
      next();
    }
  }
);

app.use(handler);

const {
  CERT_FILE = "",
  KEY_FILE = "",
  PORT = 3000,
  PRODUCTION = "FALSE",
} = config() ?? {};

if (PRODUCTION && PRODUCTION === "TRUE") {
  const privateKey = fs.readFileSync(KEY_FILE);
  const certificate = fs.readFileSync(CERT_FILE);
  const credentials = { key: privateKey, cert: certificate };

  const httpsServer = https.createServer(credentials, app);

  httpsServer.listen(PORT, () => {
    console.log("Servidor HTTPS rodando em: https://localhost:%s", PORT);
  });
} else {
  app.listen(PORT, () => {
    console.log("Servidor HTTP rodando em: http://localhost:%s", PORT);
  });
}

function config() {
  try {
    const result = dotenv.config();

    if (result.error) {
      throw result.error;
    }

    return result.parsed;
  } catch (err) {
    throw err;
  }
}
