package org.kframework.backend.go;

import edu.uci.ics.jung.graph.DirectedGraph;

import java.util.Collection;
import java.util.LinkedHashSet;
import java.util.LinkedList;
import java.util.Queue;
import java.util.Set;

// TODO: see if we still need them, is yes, move to k-distribution, or somewhere common
public class CopiedStaticMethods {

    static <V> Set<V> ancestors(
            Collection<? extends V> startNodes, DirectedGraph<V, ?> graph) {
        Queue<V> queue = new LinkedList<V>();
        Set<V> visited = new LinkedHashSet<V>();
        queue.addAll(startNodes);
        visited.addAll(startNodes);
        while (!queue.isEmpty()) {
            V v = queue.poll();
            Collection<V> neighbors = graph.getPredecessors(v);
            for (V n : neighbors) {
                if (!visited.contains(n)) {
                    queue.offer(n);
                    visited.add(n);
                }
            }
        }
        return visited;
    }

}
