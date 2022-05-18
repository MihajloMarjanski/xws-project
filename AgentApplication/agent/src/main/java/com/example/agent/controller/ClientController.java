package com.example.agent.controller;

import com.example.agent.model.*;
import com.example.agent.repository.ConfirmationTokenRepository;
import com.example.agent.service.ClientService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.context.request.WebRequest;

import java.time.Period;
import java.util.Date;

@RestController
@RequestMapping(value = "/clients")
public class ClientController {

    @Autowired
    ClientService clientService;
    @Autowired
    ConfirmationTokenRepository confirmationTokenRepository;


    @PostMapping(path = "/create")
    public ResponseEntity<?> create(@RequestBody Client client) {
        return clientService.create(client);
    }

    @GetMapping(path = "/activate")
    public ResponseEntity<?> activateClientAccount(WebRequest request, @RequestParam("token") String hashCode) {
        ConfirmationToken token = confirmationTokenRepository.findByConfirmationToken(hashCode);
        Long secs = (token.getCreatedDate().getTime() - new Date().getTime())/1000;
        Client verificationClient = token.getClient();
        if (verificationClient == null || verificationClient.isActivated() || secs > 3600 ) {
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        }

        verificationClient.setActivated(true);
        clientService.save(verificationClient);

        HttpHeaders headers = new HttpHeaders();
        headers.add("Location", "http://localhost:4200");
        return new ResponseEntity<String>(headers, HttpStatus.OK);
    }


}
