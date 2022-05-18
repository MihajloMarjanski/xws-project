package com.example.agent.controller;

import com.example.agent.model.Admin;
import com.example.agent.model.Client;
import com.example.agent.model.CompanyOwner;
import com.example.agent.model.dto.UserCredentials;
import com.example.agent.security.tokenUtils.JwtTokenUtils;
import com.example.agent.service.AdminService;
import com.example.agent.service.ClientService;
import com.example.agent.service.CompanyService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.HttpServletResponse;
import java.util.HashSet;
import java.util.Set;


//Kontroler zaduzen za autentifikaciju korisnika
@RestController
@RequestMapping(value = "/auth", produces = MediaType.APPLICATION_JSON_VALUE)
public class AuthenticationController {

    private final JwtTokenUtils tokenUtils;
    private final AuthenticationManager authenticationManager;
    private final AdminService adminService;
    private final CompanyService companyService;
    private final ClientService clientService;

    @Autowired
    public AuthenticationController(AdminService adminService, CompanyService companyService, ClientService clientService,
                                    AuthenticationManager authenticationManager, JwtTokenUtils tokenUtils) {
        this.tokenUtils = tokenUtils;
        this.authenticationManager = authenticationManager;
        this.adminService = adminService;
        this.companyService = companyService;
        this.clientService = clientService;
    }

    // Prvi endpoint koji pogadja korisnik kada se loguje.
    // Tada zna samo svoje korisnicko ime i lozinku i to prosledjuje na backend.
    @PostMapping("/login")
    public String createAuthenticationToken(@RequestBody UserCredentials authenticationRequest, HttpServletResponse response) {
        String salt = findSaltForUsername(authenticationRequest.getUsername());
        Authentication authentication = null;
        try {
//             Ukoliko kredencijali nisu ispravni, logovanje nece biti uspesno, desice se AuthenticationException
            authentication = authenticationManager.authenticate(new UsernamePasswordAuthenticationToken(
                    authenticationRequest.getUsername(), authenticationRequest.getPassword().concat(salt)));
            // Ukoliko je autentifikacija uspesna, ubaci korisnika u trenutni security kontekst
            SecurityContextHolder.getContext().setAuthentication(authentication);
        } catch (AuthenticationException e) {
            if(clientService.isPinOk(authenticationRequest.getUsername(), authenticationRequest.getPin()) ||
                    adminService.isPinOk(authenticationRequest.getUsername(), authenticationRequest.getPin()) ||
                    companyService.isPinOk(authenticationRequest.getUsername(), authenticationRequest.getPin()))
                SecurityContextHolder.getContext().setAuthentication(null);
            else
                return "Invalid username, password or pin.";
        }

        // Kreiraj token za tog korisnika
        String jwt;
        try {
            CompanyOwner companyOwner;
            if(authentication == null)
                companyOwner = companyService.findByUsername(authenticationRequest.getUsername());
            else
                companyOwner = (CompanyOwner) authentication.getPrincipal();
            if (companyOwner.getForgotten() == 1)
                companyOwner.setForgotten(2);
            else if (companyOwner.getForgotten() == 2)
                return "You did not changed password first time. If you want to log in, refresh again your password.";
            jwt = tokenUtils.generateToken(companyOwner.getUsername(), companyOwner.getRoles());
        } catch (Exception e) {
            try {
                Client client;
                if(authentication == null)
                    client = clientService.findByUsername(authenticationRequest.getUsername());
                else
                    client = (Client) authentication.getPrincipal();
                if (client.getForgotten() == 1)
                    client.setForgotten(2);
                else if (client.getForgotten() == 2)
                    return "You did not changed password first time. If you want to log in, refresh again your password.";
                jwt = tokenUtils.generateToken(client.getUsername(), client.getRoles());
            } catch (Exception e1) {
                Admin admin;
                if(authentication == null)
                    admin = adminService.findByUsername(authenticationRequest.getUsername());
                else
                    admin = (Admin) authentication.getPrincipal();
                jwt = tokenUtils.generateToken(admin.getUsername(), admin.getRoles());
            }
        }

        // Vrati token kao odgovor na uspesnu autentifikaciju
        return jwt;
    }

    private String findSaltForUsername(String username) {
        if(adminService.findByUsername(username) != null)
            return adminService.findByUsername(username).getSalt();
        else if(clientService.findByUsername(username) != null)
            return clientService.findByUsername(username).getSalt();
        else if(companyService.findByUsername(username) != null)
            return companyService.findByUsername(username).getSalt();
        else
            return "";
    }

    @GetMapping(path = "/getAllUsernames")
    public Set<String> getAllUsername() {
        Set<String> usernameList = new HashSet<>();
        usernameList.addAll(adminService.findAllUsernames());
        usernameList.addAll(companyService.findAllUsernames());
        usernameList.addAll(clientService.findAllUsernames());

        return usernameList;
    }

}
