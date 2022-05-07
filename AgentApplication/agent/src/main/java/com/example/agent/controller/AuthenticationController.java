package com.example.agent.controller;

import com.springboot.app.model.BoatOwner;
import com.springboot.app.model.Customer;
import com.springboot.app.model.Instructor;
import com.springboot.app.model.WeekendHouseOwner;
import com.springboot.app.model.dto.UserCredentials;
import com.springboot.app.security.tokenUtils.JwtTokenUtils;
import com.springboot.app.service.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.HttpServletResponse;
import java.util.Date;
import java.util.HashSet;
import java.util.Set;


//Kontroler zaduzen za autentifikaciju korisnika
@RestController
@RequestMapping(value = "/auth", produces = MediaType.APPLICATION_JSON_VALUE)
public class AuthenticationController {

    private final JwtTokenUtils tokenUtils;
    private final AuthenticationManager authenticationManager;
    private final CustomerService customerService;
    private final WeekendHouseOwnerService weekendHouseOwnerService;
    private final BoatOwnerService boatOwnerService;
    private final InstructorService instructorService;

    @Autowired
    public AuthenticationController(CustomerService customerService, WeekendHouseOwnerService weekendHouseOwnerService, BoatOwnerService boatOwnerService,
                                    AuthenticationManager authenticationManager, JwtTokenUtils tokenUtils, InstructorService instructorService) {
        this.tokenUtils = tokenUtils;
        this.authenticationManager = authenticationManager;
        this.customerService = customerService;
        this.weekendHouseOwnerService = weekendHouseOwnerService;
        this.boatOwnerService = boatOwnerService;
        this.instructorService = instructorService;
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
            Customer customer = (Customer) authentication.getPrincipal();
            if(customer.getPenalsResetingDate().getMonth() < (new Date().getMonth()) || customer.getPenalsResetingDate().getYear() <= (new Date().getYear())) {
                customer.setPenalsResetingDate(new Date());
                customer.setPenals(0);
                customerService.saveCustomer(customer);
            }
            jwt = tokenUtils.generateToken(customer.getUsername(), customer.getRole());
        } catch (Exception e) {
            try {
                BoatOwner boatOwner = (BoatOwner) authentication.getPrincipal();
                jwt = tokenUtils.generateToken(boatOwner.getUsername(), boatOwner.getRole());
            } catch (Exception e1) {
                try {
                    Instructor instructor = (Instructor) authentication.getPrincipal();
                    jwt = tokenUtils.generateToken(instructor.getUsername(), instructor.getRole());
                } catch (Exception e2) {
                    WeekendHouseOwner weekendHouseOwner = (WeekendHouseOwner) authentication.getPrincipal();
                    jwt = tokenUtils.generateToken(weekendHouseOwner.getUsername(), weekendHouseOwner.getRole());
                }
            }
        }

        // Vrati token kao odgovor na uspesnu autentifikaciju
        return jwt;
    }

    @GetMapping(path = "/getAllUsernames")
    public Set<String> getAllUsername() {
        Set<String> usernameList = new HashSet<String>();
        usernameList.addAll(customerService.findAllUsernames());
        usernameList.addAll(weekendHouseOwnerService.findAllUsernames());
        usernameList.addAll(boatOwnerService.findAllUsernames());
        usernameList.addAll(instructorService.findAllUsernames());

        return usernameList;
    }

}
