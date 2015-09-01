defmodule Horoscope do
  use Application
  require Logger

  def worker_pool, do: :horoscope_workers
  def worker_pool_config do
    [worker_module: Horoscope.Worker,
     name: {:local, worker_pool},
     size: 32,
     max_overflow: 8]
  end

  def start(_type, _args) do
    import Supervisor.Spec

    children = [
      worker(__MODULE__, [], function: :server),
      worker(Horoscope.Repo, []),
      :poolboy.child_spec(worker_pool, worker_pool_config, [])
    ]

    options = [
      strategy: :one_for_one,
      name: Algorithm.Supervisor
    ]

    Supervisor.start_link(children, options)
  end

  def server do
    Logger.info("Server running on port 4000")
    {:ok, _} = Plug.Adapters.Cowboy.http(Horoscope.Router, [])
  end
end
