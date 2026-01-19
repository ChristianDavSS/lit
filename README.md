# **lit: Your project best friend 🥇**
### **What's 'lit'? ❓**
CLI Tool developed 100% in Golang in order to help you get quick information about your git respositories and project metrics. Thus, it can be really helpful for big work enviroments in order to improve productivity or detect code anomalies,
getting warnings about possible dangerous code.

### **Top features 🔝**
With the available commands, you can analyze anything you need from code to even repository metrics. The scanning does a search in your working directory looking for the scripts written on the supported languages to scan and calculates the cyclical
complexity so you can get feedback based in the metric.

### **Future features (plans)**
I'll be implementing the scanner for classes/structs very soon, as well as the conventions for the naming of the functions/methods.

### **Supported languages for the code scan 💻**
Current supported languages: **Java, Python,  JavaScript, Go**

Incoming support for the next languages: **C, C++, C#, TypeScript**

### **Requirements 🧰**
- 64x C compiler installed (for the file scanner)
  
  That's the only requirement to run the project.

### **Available commands 🌝**
- *lit files*: The brain and main command of this project. This command itself scans your whole repository and finds the scripts of the supported languages, scanning and looking for possible dangerous functions defined
  and variables that doesn't match with the naming convention defined.
  *Command flags*:
  - *--loc*: This flag allows you to know how much lines of code were written on which language and their percentage of the total lines.

- *lit config*: With this, you can modify the **config.json** file with a friendly and easy interface. Currently, the only configuration is the current regex for the naming convention you're using.
  I would recommend to run this command before anything. (the default convention is **CamelCase/camelCase**).

  The available conventions are: **camelCase, CamelCase, snake_case, CamelCase/camelCase** but if you want to use your own you just have to modify the value of the **'activeNamingConvention'** key on the .json file.

### Why was this developed?
I developed ***Lit*** because I wanted to reinforce my Go knowledge by creating a useful and meaningful project that developers like me could use in their best projects to keep the code cleaner. Writting clean but efficient code is very important
because the code will always be read by developers and they must understand it. Developing this project I learned about go routines and how powerful they are, also, I was able to reinforce my Go knowledge and now I feel ready to start bigger projects.

### Support the project ♥️
Your **star ⭐** in the repository would be more than enough.

**Project open for contributions. Pull-requests are welcome :D.**
