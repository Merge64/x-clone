---
name: Non-Functional Requirement
about: Requirement that's not a feature itself but needs to be met
title: NFR - Current Requirement Name
labels: NFR
assignees: ''

---

### **Non-Functional Requirement (NFR) Template**

#### **1. Title**  
   - A concise and descriptive title for the NFR.  
     Example: **System Performance Under Load**

---
#### **3. Description**  
   - A detailed description of what the NFR entails.  
     Example:  
     _The system should maintain a response time of under 2 seconds for 95% of all user requests under a load of 1,000 concurrent users._

---

#### **3. Category**  
   - Specify the category of the NFR (choose one or add a custom category):  
     - **Performance**  
     - **Security**  
     - **Scalability**  
     - **Usability**  
     - **Availability**  
     - **Compliance**  
     - **Maintainability**  
     - **Portability**  

---

#### **4. Acceptance Criteria**  
   - Define how this NFR will be measured and validated.  
     Example:  
     - Response time logs will be collected during a stress test simulating 1,000 concurrent users.  
     - 95% of requests should complete in under 2 seconds.
---

#### **5. Priority**  
   - Indicate the importance of this requirement.  
     Example:  
     - High  
     - Medium  
     - Low  

---

#### **6. Dependencies**  
   - List any other requirements, systems, or tasks that this NFR depends on.  
     Example:  
     _Dependent on database optimization (Task-DB-004)._
