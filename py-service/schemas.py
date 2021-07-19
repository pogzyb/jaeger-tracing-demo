from typing import List, Any

from pydantic import BaseModel


class MatrixOut(BaseModel):
    rows: int
    columns: int
    # zeroes: bool
    # ones: bool
    matrix: List[Any]  # is there a typing for multi-dimensional arrays?
