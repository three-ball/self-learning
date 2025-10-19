# STACK

![Stack](../images/stack.png)

## Definition

- **A Stack is a linear data structure that follows the LIFO (Last In, First Out) principle**. Think of it like a stack of plates - you can only add or remove plates from the top. The last element added to the stack will be the first one to be removed.

## Characteristics

- **LIFO Principle**: Last In, First Out
- **Top Pointer**: Tracks the top element
- **Two Primary Operations**: `Push` (add) and `Pop` (remove)
- **Time Complexity**: O(1) for push, pop, and peek operations


```mermaid
graph TD
    A[Stack Operations] --> B[Push Operation]
    A --> C[Pop Operation]
    A --> D[Peek Operation]
    
    B --> B1[Check if Full]
    B1 -->|Not Full| B2[Increment Top]
    B2 --> B3[Add Element at Top]
    B1 -->|Full| B4[Stack Overflow Error]
    
    C --> C1[Check if Empty]
    C1 -->|Not Empty| C2[Get Element at Top]
    C2 --> C3[Decrement Top]
    C1 -->|Empty| C4[Stack Underflow Error]
    
    D --> D1[Check if Empty]
    D1 -->|Not Empty| D2[Return Top Element]
    D1 -->|Empty| D3[Return Error]
    
    style A fill:#f9f,stroke:#333,stroke-width:4px
    style B fill:#bbf,stroke:#333,stroke-width:2px
    style C fill:#bfb,stroke:#333,stroke-width:2px
    style D fill:#fbb,stroke:#333,stroke-width:2px
```

## Real-world Applications

- Browser History Management
- Undo/Redo Functionality in Text Editors