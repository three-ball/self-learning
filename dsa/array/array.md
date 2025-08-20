# ARRAY

## Definition

An **Array** is a collection of elements stored in contiguous memory locations, where each element can be accessed using an index. Arrays are homogeneous data structures, meaning all elements are of the same data type.

## Characteristics

- **Fixed Size**: In most languages (including Go), arrays have a fixed size determined at compile time
- **Contiguous Memory**: Elements are stored in contiguous memory locations, allowing for efficient access
- **Homogeneous**: All elements in an array are of the same type
- **Index-Based Access**: Elements can be accessed using an index, **starting from 0**.

## Key Learning Points

- Memory Efficiency: Arrays provide excellent cache locality due to contiguous memory allocation
- Trade-offs: Fast access (O(1)) vs. expensive insertions/deletions (O(n))
- Go Specifics:
    - Arrays have fixed size: `[5]int`.
    - Slice are dynamic arrays: `[]int`.
    - Use slices for dynamic array needs in Go, as they provide more flexibility than fixed-size arrays.