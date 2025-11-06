import pandas as pd
import matplotlib.pyplot as plt
import japanize_matplotlib


def drop_duplicates():
    df = pd.read_csv("../copy_251031.csv", header=None, names=["time", "kind"])
    df = df.drop_duplicates()

    df.to_csv("modified_251031.csv", index=False)


def init_figure(file_path: str, graph_title: str):
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


init_figure("modified_251031.csv", "30分ごとの種類別来場者数 10/31")
init_figure("../log_251101.csv", "30分ごとの種類別来場者数 11/1")
init_figure("../log_251102.csv", "30分ごとの種類別来場者数 11/2")
plt.show()
