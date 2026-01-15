# **lit: Your project best friend 🥇**
### **What's 'lit'? ❓**
CLI Tool developed 100% in Golang in order to help you get quick information about your git respositories and project metrics. Thus, it can be really helpful for big work enviroments in order to improve productivity or detect code anomalies,
getting warnings about possible dangerous code.

### **Top features 🔝**
With the available commands, you can analyze anything you need from code to even repository metrics. The scanning does a search in your working directory looking for the scripts written on the supported languages to scan and calculates the cyclical
complexity so you can get feedback based in the metric.

### **Future features (plans)**
I want to add more specific filters to the scanner so you can scan the code by your favorite code conventions (***camelCase, CamelCase, sneak_case, etc***).

### **Supported languages for the code scan 💻**
Current supported languages: **Java, Python,  JavaScript**

Incoming support for the next languages: **C, C++, C#, Go, TypeScript**

### **Requirements 🧰**
- 64x C compiler installed (for the file scanner)
  
  That's the only requirement to run this project :D

### **Available commands 🌝**
- *lit authors*: Returns the authors of the project with their total commits. ***This command and it's flags are still being developed***
  
  *Command flags*:
  - *--verbose*: This will show every commit information (hash, date and message, just like a **git log** but sorted by authors)
  - *--commit-size*: Print out the additions, deletions and total lines of code modified on that specific commit
  - *-w, --who*: With this flag, you can search a specific user commits by their username or email (case sensitive).
  - *--stats*: This command shows the stadistics of every file modified in a specific commit.
  - *--since*: Takes a date in the DD/MM/YYYY format and prints out the commits from that date on.
  - *--until*: Takes a date in the DD/MM/YYYY format and prints out the commits until that date.
 
- *lit files*: The brain and main command of this project. This command itself scans your whole repository and finds the scripts of the supported languages, scanning and looking for possible dangerous functions defined.
  *Command flags*:
  - *--loc*: This flag allows you to know how much lines of code were written on which language and their percentage of the total lines.

### Why was this developed?
I developed ***Lit*** because I wanted to reinforce my Go knowledge by creating a useful and meaningful project that developers like me could use in their best projects to keep the code cleaner. Writting clean but efficient code is very important
because the code will always be read by developers and they must understand it. Developing this project I learned about go routines and how powerful they are, also, I was able to reinforce my Go knowledge and now I feel ready to start bigger projects.

### Support the project ♥️
Your **star ⭐** in the repository would be more than enough.

**Project open for contributions. Pull-requests are welcome :D.**
