import pandas as pd
import matplotlib.pyplot as plt


def drop_duplicates():
    df = pd.read_csv("../copy_251031.csv", header=None, names=["time", "kind"])
    df = df.drop_duplicates()

    df.to_csv("modified_251031.csv", index=False)


def show_graph():
    df = pd.read_csv("modified_251031.csv")
    df_grouped = (
        df.groupby([pd.Grouper(key="time", freq="30 min"), "kind"])
        .size()
        .unstack(fill_value=0)
    )
    df_grouped.plot(
        kind="bar",
        stacked=True,
        colormap="tab10",  # 色の指定（例："tab10", "Set2", "Paired" など）
        figsize=(10, 5),
    )
    plt.title("Count of 'kind' every 10 minutes")
    plt.xlabel("Time (10min bins)")
    plt.ylabel("Count")
    plt.xticks(rotation=45)
    plt.tight_layout()
    plt.show()


show_graph()
