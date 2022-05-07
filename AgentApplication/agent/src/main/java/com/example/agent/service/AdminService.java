package com.example.agent.service;

import com.example.agent.repository.AdminRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class AdminService {

    @Autowired
    private RoleService roleService;
    @Autowired
    private AdminRepository adminRepository;
}
