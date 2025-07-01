package main

import (
	"fmt"
)

func main() {
	manager := Manager{}
	manager.AddEmployee(Employee{ID: 1, Name: "wahono", Age: 1, Salary: 80000})
	manager.AddEmployee(Employee{ID: 2, Name: "ridho", Age: 12, Salary: 40000})
	resultAVG := manager.GetAverageSalary()
	manager.RemoveEmployee(1)
	resultFind := manager.FindEmployeeByID(2)
	fmt.Println(resultFind, resultAVG)
}

func (m *Manager) AddEmployee(e Employee) {
	m.Employees = append(m.Employees, e)
}

func (m *Manager) RemoveEmployee(id int) {
	for i, r := range m.Employees {
		if r.ID == id {
			m.Employees[i] = m.Employees[len(m.Employees)-1]
			m.Employees = m.Employees[:len(m.Employees)-1]
			return
		}
	}
	return
}

func (m *Manager) GetAverageSalary() float64 {
	length := len(m.Employees)
	var TotalSalaries float64 = 0
	for _, r := range m.Employees {
		TotalSalaries += r.Salary
	}
	AVGSalaries := TotalSalaries / float64(length)
	return AVGSalaries
}

func (m *Manager) FindEmployeeByID(id int) *Employee {
	for _, r := range m.Employees {
		if r.ID == id {
			return &r
		}
	}
	return nil
}

type Employee struct {
	ID     int
	Name   string
	Age    int
	Salary float64
}

type Manager struct {
	Employees []Employee
}
