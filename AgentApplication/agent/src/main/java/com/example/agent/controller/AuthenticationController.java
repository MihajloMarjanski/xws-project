package com.example.agent.controller;

import com.example.agent.model.Admin;
import com.example.agent.model.CompanyOwner;
import com.example.agent.model.dto.UserCredentials;
import com.example.agent.security.tokenUtils.JwtTokenUtils;
import com.example.agent.service.AdminService;
import com.example.agent.service.CompanyService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.HttpServletResponse;


//Kontroler zaduzen za autentifikaciju korisnika
@RestController
@RequestMapping(value = "/auth", produces = MediaType.APPLICATION_JSON_VALUE)
public class AuthenticationController {

    private final JwtTokenUtils tokenUtils;
    private final AuthenticationManager authenticationManager;
    private final AdminService adminService;
    private final CompanyService companyService;

    @Autowired
    public AuthenticationController(AdminService adminService, CompanyService companyService,
                                    AuthenticationManager authenticationManager, JwtTokenUtils tokenUtils) {
        this.tokenUtils = tokenUtils;
        this.authenticationManager = authenticationManager;
        this.adminService = adminService;
        this.companyService = companyService;
    }

    // Prvi endpoint koji pogadja korisnik kada se loguje.
    // Tada zna samo svoje korisnicko ime i lozinku i to prosledjuje na backend.
    @PostMapping("/login")
    public String createAuthenticationToken(@RequestBody UserCredentials authenticationRequest, HttpServletResponse response) {

        // Ukoliko kredencijali nisu ispravni, logovanje nece biti uspesno, desice se AuthenticationException
        Authentication authentication = authenticationManager.authenticate(new UsernamePasswordAuthenticationToken(
                authenticationRequest.getUsername(), authenticationRequest.getPassword()));

        // Ukoliko je autentifikacija uspesna, ubaci korisnika u trenutni security kontekst
        SecurityContextHolder.getContext().setAuthentication(authentication);

        // Kreiraj token za tog korisnika
        String jwt;
        try {
            CompanyOwner companyOwner = (CompanyOwner) authentication.getPrincipal();
            jwt = tokenUtils.generateToken(companyOwner.getUsername(), companyOwner.getRole());
        } catch (Exception e) {
            Admin admin = (Admin) authentication.getPrincipal();
            jwt = tokenUtils.generateToken(admin.getUsername(), admin.getRole());
        }

        // Vrati token kao odgovor na uspesnu autentifikaciju
        return jwt;
    }

//    @GetMapping(path = "/getAllUsernames")
//    public Set<String> getAllUsername() {
//        Set<String> usernameList = new HashSet<String>();
//        usernameList.addAll(customerService.findAllUsernames());
//        usernameList.addAll(weekendHouseOwnerService.findAllUsernames());
//        usernameList.addAll(boatOwnerService.findAllUsernames());
//        usernameList.addAll(instructorService.findAllUsernames());
//
//        return usernameList;
//    }

}
