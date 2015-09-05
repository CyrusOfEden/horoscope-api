defmodule Horoscope do
  use Application
  require Logger

  def worker_pool, do: :horoscope_workers
  def worker_pool_config do
    [worker_module: Horoscope.Worker,
     name: {:local, worker_pool},
     size: 4,
     max_overflow: 2]
  end

  def start(_type, _args) do
    import Supervisor.Spec

    children = [
      worker(Horoscope.Repo, []),
      :poolboy.child_spec(worker_pool, worker_pool_config, [])
    ]

    options = [
      strategy: :one_for_one,
      name: Algorithm.Supervisor
    ]

    Supervisor.start_link(children, options)
  end

  defdelegate get, to: Horoscope.Worker, as: :call
  defdelegate get(params), to: Horoscope.Worker, as: :call
end
