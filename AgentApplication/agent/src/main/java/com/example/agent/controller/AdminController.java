package com.example.agent.controller;

import com.example.agent.model.Admin;
import com.example.agent.model.Client;
import com.example.agent.model.CompanyOwner;
import com.example.agent.model.Role;
import com.example.agent.service.AdminService;
import com.example.agent.service.RoleService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.HttpServletRequest;
import java.util.List;

@RestController
@RequestMapping(value = "/admin")
public class AdminController {

    @Autowired
    private AdminService adminService;


    @PreAuthorize("hasRole('ADMIN')")
    @PutMapping(path = "/approve/company/{id}")
    public ResponseEntity<?> approveCompany(@PathVariable Integer id, HttpServletRequest request) {
        return adminService.approveCompany(id, request);
    }

    @PreAuthorize("hasRole('ADMIN')")
    @GetMapping(path = "/username/{username}")
    public ResponseEntity<?> adminByUsername(@PathVariable String username , HttpServletRequest request) {
        return adminService.getByUsername(username, request);
    }

    @PreAuthorize("hasRole('ADMIN')")
    @PostMapping(path = "/update")
    public ResponseEntity<?> updateAdmin(@RequestBody Admin client, HttpServletRequest request) {
        return adminService.updateAdmin(client , request);
    }
}
