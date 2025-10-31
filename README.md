# Graphographic

A simple graph visualizer with three algorithms build in (DFS, BFS and Dijkstra shortest path between two nodes).

## Controls

The program is based on modes that you can switch between with certain key presses.

Use the right mouse button to move around and BACKSPACE to revert actions.

### Place mode

Enabled with the P key, lets you place new nodes with a left click.

### Move mode

Enabled with the M key, lets you move nodes.

### Edit mode

Enabled with the E key, lets you edit node names and edge lengths.

### Append mode

Enabled with the A key, lets you create by extending it from an already existing one. Left click and hold a node and move the new node somewhere.

### Connect mode

Enabled with the C key, lets you connect two nodes. SHIFT+D changes the connecting behavior:

- DIRECTED -- create a one way connecting edge from node A to B
- UNDIRECTED -- create a two way connecting edge from node A to B (it creates two edges)

### Algorithm mode

Enabled with the T key, lets you execute build-in algorithms on created graphs. Algorithms expect one or more nodes to be selected and report errors if those requirements are not met. You can execute an algorithm with the R key.

### Delete mode

Enabled with the D key, lets you delete nodes and edges.
