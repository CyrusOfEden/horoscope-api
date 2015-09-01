defmodule Horoscope.Router do
  use Plug.Router
  alias Horoscope.Worker

  Mix.env == :dev && use Plug.Debugger

  plug :match
  plug :dispatch

  def send_json(conn, code, nil) do
    send_json(conn, code, %{error: "invalid request"})
  end
  def send_json(conn, code, response) do
    conn
    |> put_resp_content_type("application/json")
    |> send_resp(code, Poison.encode!(response))
  end

  get "/" do
    send_json(conn, 200, Worker.call)
  end

  get "/:sign" do
    send_json(conn, 200, Worker.call(sign))
  end

  get "/:year/:week" do
    send_json(conn, 200, Worker.call({year, week}))
  end

  get "/:year/:week/:sign" do
    send_json(conn, 200, Worker.call({year, week, sign}))
  end

  match _ do
    send_json(conn, 400, nil)
  end
end