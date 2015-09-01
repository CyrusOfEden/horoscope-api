defmodule Horoscope.Router do
  use Plug.Router
  alias Horoscope.Worker

  Mix.env == :dev && use Plug.Debugger

  plug :match
  plug :dispatch

  @error Poison.encode!(%{error: "invalid request"})
  def send_json(conn, code, nil) do
    send_json(conn, code, @error)
  end
  def send_json(conn, code, response) do
    conn
    |> put_resp_content_type("application/json")
    |> send_resp(code, response)
  end

  def horoscope(conn, params) do
    case Worker.call(Map.put(params, :encode, true)) do
      "null"   -> send_json(conn, 400, nil)
      response -> send_json(conn, 200, response)
    end
  end

  get "/",                  do: horoscope(conn, %{})
  get "/:sign",             do: horoscope(conn, %{sign: sign})
  get "/:year/:week",       do: horoscope(conn, %{week: {year, week}})
  get "/:year/:week/:sign", do: horoscope(conn, %{week: {year, week}, sign: sign})

  match _, do: send_json(conn, 400, nil)
end