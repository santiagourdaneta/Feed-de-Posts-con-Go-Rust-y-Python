use actix::prelude::*;
use actix_web::{web, App, HttpResponse, HttpServer, Responder, HttpRequest, Error};
use actix_web_actors::ws; // This is the line you were missing!

// MyWs is the actor that will handle the WebSocket connection
struct MyWs;

impl Actor for MyWs {
    type Context = ws::WebsocketContext<Self>;
}

impl StreamHandler<Result<ws::Message, ws::ProtocolError>> for MyWs {
    fn handle(&mut self, msg: Result<ws::Message, ws::ProtocolError>, ctx: &mut Self::Context) {
        // This is where you handle incoming messages
        match msg {
            Ok(ws::Message::Ping(msg)) => ctx.pong(&msg),
            Ok(ws::Message::Text(text)) => {
                println!("Received text: {:?}", text);
                ctx.text(format!("Hello! You sent: {}", text));
            }
            Ok(ws::Message::Close(reason)) => {
                ctx.close(reason);
                ctx.stop();
            }
            _ => ctx.stop(),
        }
    }
}

async fn index() -> impl Responder {
    HttpResponse::Ok().body("Hola, soy el servidor de notificaciones!")
}

async fn ws_route(req: HttpRequest, stream: web::Payload) -> Result<HttpResponse, Error> {
    ws::start(MyWs, &req, stream)
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .route("/", web::get().to(index))
            .route("/ws/", web::get().to(ws_route))
    })
    .bind("127.0.0.1:8081")?
    .run()
    .await
}