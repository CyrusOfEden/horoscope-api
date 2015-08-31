defmodule Horoscope do
  use Application

  def start(_type, _args) do
    import Supervisor.Spec

    children = [
      worker(Horoscope.Repo, []),
      worker(Horoscope.Worker, [])
    ]

    opts = [strategy: :one_for_one, name: Algorithm.Supervisor]
    Supervisor.start_link(children, opts)
  end
end
