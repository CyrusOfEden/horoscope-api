defmodule Horoscope do
  use Application

  def start(_type, _args) do
    import Supervisor.Spec

    children = [
      # worker(Horoscope.Worker, [arg1, arg2, arg3])
    ]

    opts = [strategy: :one_for_one, name: Horoscope.Supervisor]
    Supervisor.start_link(children, opts)
  end
end
