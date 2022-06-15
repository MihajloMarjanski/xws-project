package com.example.agent.controller;

import com.example.agent.model.Client;
import com.example.agent.model.dto.ActivationLinkDto;
import com.example.agent.service.ClientService;
import com.example.agent.service.EmailService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping(value = "/email")
public class EmailController {
    @Autowired
    EmailService emailService;


    @PostMapping(path = "/activation")
    public ResponseEntity<?> activationLink(@RequestBody ActivationLinkDto dto) {
        if (!dto.getApiKey().equals(""))
            return emailService.sendActivationMailToDislinktUser(dto.getEmail(), dto.getName(), dto.getApiKey());
        emailService.sendNewPassword(dto.getEmail(), dto.getName());
        return new ResponseEntity<>(HttpStatus.OK);
    }
}
