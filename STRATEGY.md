# Sampling Strategy for Data Availability Sampling

## Background

The data availability sampling strategy is designed to ensure that data remains accessible on the network. This approach involves sampling portions of the data across the network to verify its presence and integrity. Implemented by the light-client, this strategy helps prevent data loss or corruption. Erasure coding plays a crucial role in this process by encoding data into redundant fragments, which are then distributed across the network. The way the data is erasure coded directly impacts the effectiveness of the sampling strategy, as it determines how easily the light-client can verify data availability. By sampling and verifying these fragments, the light-client ensures the data’s reliability and accessibility on the network.

## 1D Erasure Coding

In this system, 1D erasure coding is applied to data organized into rows, where each row is extended to 128 cells. The first 64 cells contain the original data, while the remaining 64 cells store redundancy generated through the erasure coding process. This structure allows for the recovery of the original data even if some cells are lost or corrupted. However, to achieve a high level of confidence in data integrity—such as 99%—it is crucial to test every row. By verifying each row individually, we ensure that the erasure coding effectively protects the data and that the system maintains its reliability across all rows.

### Samples number estimation

$$
S \geq R - R \times \left(1 - C\right)^{\frac{1}{E}}
$$

Where:

- $S$ is the number of samples needed.
- $R$ is the total number of cells per row (128 in this case).
- $C$ is the desired confidence level.
- $E$ is the number of cells needed to recover the entire row (64 in this case).

How It Works:

- **Confidence Level** $C$: This represents the desired probability that the row can be correctly reconstructed using $E$ intact cells.
- **Number of Samples** $S$: The formula calculates the minimum number of samples needed from each row to achieve the desired confidence level.

Example:

Let’s calculate the number of samples $S$ needed for this specific case:

- $R = 128$ (total number of cells per row)
- $C = 0.99$ (99% confidence level)
- $E = 64$ (number of cells needed to recover the entire row)

$$
S \geq 128 - 128 \times \left(1 - 0.99\right)^{\frac{1}{64}}
$$

$$
S \geq 128 - 128 \times 0.9934427 \approx 8.89
$$

This result indicates that you would need to sample at least 8.89 cells per row to achieve a 99% confidence level that at least 64 of the 128 cells are intact, ensuring that the row can be fully recovered.

If to put the results in percentage perspective, you would need to sample at least 7% of the cells in each row to achieve a 99% confidence level that the row can be fully recovered. The value remains the same for all rows, as the erasure coding process is consistent across the data.