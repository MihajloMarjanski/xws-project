package com.example.agent.controller;

import com.example.agent.model.dto.ActivationLinkDto;
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
        if (dto.getApiKey().equals("")) {
            emailService.sendNewPassword(dto.getEmail(), dto.getName());
            return new ResponseEntity<>(HttpStatus.OK);
        }
        else if (dto.getName().equals("")) {
            emailService.send2factorAuthPin(dto.getEmail(), dto.getApiKey());
            return new ResponseEntity<>(HttpStatus.OK);
        }
        else if (dto.getName().equals("token")) {
            emailService.sendPasswordlesstoDislinkt(dto.getEmail(), dto.getApiKey());
            return new ResponseEntity<>(HttpStatus.OK);
        }
        return emailService.sendActivationMailToDislinktUser(dto.getEmail(), dto.getName(), dto.getApiKey());
    }
}
