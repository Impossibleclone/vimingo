package main


func isInSelection(start, end Cursor, x, y int) bool {
    if start.Y > end.Y || (start.Y == end.Y && start.X > end.X) {
        start, end = end, start
    }

    if y < start.Y || y > end.Y {
        return false
    }

    if y == start.Y && y == end.Y {
        return x >= start.X && x < end.X
    } else if y == start.Y {
        return x >= start.X
    } else if y == end.Y {
        return x < end.X
    }

    return true
}

