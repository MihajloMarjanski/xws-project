package com.example.agent.controller;

import com.example.agent.service.ClientService;
import com.example.agent.service.EmailService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(value = "/email")
public class EmailController {
    @Autowired
    EmailService emailService;


    @GetMapping(path = "/activation/{email}/{name}/{key}")
    public ResponseEntity<?> activateUser(@PathVariable String email, @PathVariable String name, @PathVariable String key) {
        return emailService.sendActivationMailToDislinktUser(email, name, key);
    }
}
