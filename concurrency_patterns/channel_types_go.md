## 📚 Channel Types in Go

There are **THREE types** of channels:

---

### **1. Bidirectional Channel (can send AND receive)**

```go
ch := make(chan int)

ch <- 5      // ✅ Can send
val := <-ch  // ✅ Can receive
```

**Type:** `chan int`

---

### **2. Send-Only Channel (can only send)**

```go
var ch chan<- int  // Notice the arrow AFTER chan

ch <- 5      // ✅ Can send
val := <-ch  // ❌ COMPILE ERROR! Can't receive
```

**Type:** `chan<- int` (arrow points INTO channel)

---

### **3. Receive-Only Channel (can only receive)**

```go
var ch <-chan int  // Notice the arrow BEFORE chan

ch <- 5      // ❌ COMPILE ERROR! Can't send
val := <-ch  // ✅ Can receive
```

**Type:** `<-chan int` (arrow points OUT OF channel)

---

## 🎯 Visual Memory Aid:

```go
chan int      // ↔️ Bidirectional (no arrow, can do both)

chan<- int    // ➡️ Send-only (arrow INTO channel)
              //    Think: "I can only push INTO it"

<-chan int    // ⬅️ Receive-only (arrow OUT OF channel)
              //    Think: "I can only pull FROM it"
```

---

## 💡 Why Do We Need This?

**For safety and clarity!**

### **Example: Generator Function**

```go
func Generator(nums ...int) <-chan int {  // ⭐ Returns receive-only
    out := make(chan int)  // Create bidirectional
    
    go func() {
        for _, n := range nums {
            out <- n  // Send to it (inside function)
        }
        close(out)
    }()
    
    return out  // ⭐ Return as receive-only
}
```

**What happens:**

1. **Inside function:** Channel is bidirectional (can send)
2. **Return type:** Converts to receive-only
3. **Caller:** Can only receive, can't send (safety!)

---

## 🔥 Real Example:

```go
func Producer() <-chan int {  // Returns receive-only
    ch := make(chan int)
    go func() {
        ch <- 1
        ch <- 2
        ch <- 3
        close(ch)
    }()
    return ch
}

func main() {
    numbers := Producer()
    
    // ✅ Can receive
    val := <-numbers
    
    // ❌ Can't send (compile error!)
    numbers <- 100  // ERROR: invalid operation
}
```

**This prevents bugs!** You can't accidentally send to a channel that should only be read from!

---

## 📊 In Your Code:

### **Square Function:**

```go
func Square(in <-chan int) <-chan int {
//           ^^^^^^^^^^^   ^^^^^^^^^^^
//           Receive-only  Receive-only
//           parameter     return type
    
    out := make(chan int)  // Inside: bidirectional
    
    go func() {
        for n := range in {      // ✅ Read from in
            out <- n * n         // ✅ Write to out
        }
        close(out)
    }()
    
    return out  // Returned as receive-only
}
```

**Why?**
- `in <-chan int`: Square should only READ from input
- `<-chan int` return: Caller should only READ from output

---

## 🎯 So Your Bug:

```go
worker := make([]chan int, numWorkers)  // ❌ Slice of bidirectional
//                ^^^^^^^

for i := 0; i < numWorkers; i++ {
    worker = append(worker, Square(i, input))
    //                      ^^^^^^^^^^^^^^^^
    //                      Returns <-chan int (receive-only)
}
```

**Problem:** Trying to put receive-only channel into bidirectional slice!

---

## ✅ Fixed:

```go
workers := make([]<-chan int, numWorkers)  // ✅ Slice of receive-only
//                 ^^^^^^^^^

for i := 0; i < numWorkers; i++ {
    workers[i] = Square(i, input)  // ✅ Matches!
}
```

---

## 💪 Quick Reference:

| Syntax | Name | Can Send? | Can Receive? |
|--------|------|-----------|--------------|
| `chan int` | Bidirectional | ✅ | ✅ |
| `chan<- int` | Send-only | ✅ | ❌ |
| `<-chan int` | Receive-only | ❌ | ✅ |

---

## 🔥 Practice Understanding:

```go
// Function that PRODUCES data (returns receive-only)
func Producer() <-chan int {
    ch := make(chan int)
    go func() {
        ch <- 1  // Can send inside
    }()
    return ch  // Caller can only receive
}

// Function that CONSUMES data (takes receive-only)
func Consumer(in <-chan int) {
    val := <-in  // Can only receive
}

// Function that both reads and writes (takes bidirectional)
func Processor(ch chan int) {
    ch <- 5      // Can send
    val := <-ch  // Can receive
}
```

---

## ✅ Does This Make Sense Now?

**Key points:**
1. ✅ `chan int` = bidirectional (send AND receive)
2. ✅ `chan<- int` = send-only (arrow INTO channel)
3. ✅ `<-chan int` = receive-only (arrow OUT OF channel)
4. ✅ Use directional channels for safety and clarity




