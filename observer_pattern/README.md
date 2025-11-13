type PhoneDisplay struct { }
func (p *PhoneDisplay) Update(data) { }
```
**Purpose:** Each observer decides HOW to react to updates

---

## 📊 Diagram:
```
WeatherStation (Subject)
    |
    |--- Subscribe(phone)
    |--- Subscribe(tv)
    |--- Subscribe(window)
    |
    SetTemperature(25) ──→ NotifyAll()
                              |
                              ├──→ phone.Update(25)
                              ├──→ tv.Update(25)
                              └──→ window.Update(25)



🎓 Benefits of Observer Pattern:

✅ Loose Coupling - Weather station doesn't need to know about specific displays
✅ Easy to Add/Remove - Subscribe/unsubscribe observers anytime
✅ Automatic Updates - One change, everyone notified
✅ Flexible - Each observer can react differently