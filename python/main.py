import pandas as pd
import matplotlib.pyplot as plt
import japanize_matplotlib
from pathlib import Path
from pandas import Series


def init_bar_graph(file_path: str, graph_title: str):
    df = pd.read_csv(file_path, header=None, names=["time", "kind"])
    df["time"] = df["time"].map(lambda s: s.split(" +0900")[0])
    df["time"] = pd.to_datetime(df["time"], format="%Y-%m-%d %H:%M:%S.%f")
    df.set_index("time", inplace=True)
    df_grouped = (
        df.groupby([pd.Grouper(freq="30 min"), "kind"]).size().unstack(fill_value=0)
    )
    ax = df_grouped.plot(
        kind="bar",
        stacked=True,
        colormap="tab10",
        figsize=(10, 5),
    )
    ax.set_title(graph_title)
    ax.set_xlabel("Time")
    ax.set_ylabel("Count")


def make_autopct(values: Series[int]):
    def my_autopct(pct):
        total = sum(values)
        val = int(round(pct * total / 100.0))
        return f"{pct:.1f}%\n({val})"

    return my_autopct


def init_pie_chart(file_path: str | None, graph_title: str):
    """`file_path`が`None`のときは`logs/`の下のファイルを全部見る"""
    if file_path is None:
        folder = Path("../logs/")
        df = pd.concat(
            [
                pd.read_csv(file, header=None, names=["time", "kind"])
                for file in folder.glob("*.csv")
            ],
            ignore_index=True,
        )
    else:
        df = pd.read_csv(file_path, header=None, names=["time", "kind"])

    counts = df["kind"].value_counts()
    _, ax = plt.subplots()

    result = ax.pie(
        counts,
        startangle=90,
        counterclock=False,
        autopct=make_autopct(counts),
        pctdistance=0.8,
    )
    ax.set_title(graph_title)
    ax.legend(
        result[0],
        list(counts.index),
        title="Kinds",
        loc="center left",
        bbox_to_anchor=(1, 0, 0.5, 1),
    )


"""init_bar_graph("../logs/log_251031.csv", "30分ごとの種類別来場者数 10/31")
init_bar_graph("../logs/log_251101.csv", "30分ごとの種類別来場者数 11/1")
init_bar_graph("../logs/log_251102.csv", "30分ごとの種類別来場者数 11/2")

init_pie_chart("../logs/log_251031.csv", "来場者の割合 10/31")
init_pie_chart("../logs/log_251101.csv", "来場者の割合 11/1")
init_pie_chart("../logs/log_251102.csv", "来場者の割合 11/2")"""
init_pie_chart(None, "来場者の割合 合計")

plt.show()
