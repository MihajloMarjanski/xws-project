package com.example.agent.controller;

import com.example.agent.model.Client;
import com.example.agent.model.ConfirmationToken;
import com.example.agent.model.dto.UserDto;
import com.example.agent.repository.ConfirmationTokenRepository;
import com.example.agent.service.ClientService;
import com.example.agent.service.CompanyService;
import lombok.extern.slf4j.Slf4j;
import org.slf4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.context.request.WebRequest;

import javax.servlet.http.HttpServletRequest;
import javax.validation.Valid;
import java.util.Date;

@Slf4j
@RestController
@RequestMapping(value = "/clients")
public class ClientController {

    @Autowired
    ClientService clientService;
    @Autowired
    ConfirmationTokenRepository confirmationTokenRepository;
    @Autowired
    CompanyService companyService;


    @PostMapping(path = "/create")
    public ResponseEntity<?> create(@Valid @RequestBody UserDto client, BindingResult res, HttpServletRequest request) {
        if(res.hasErrors()){
            log.warn("Ip: {}, Fields for new client not valid!", request.getRemoteAddr());
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        }
        return clientService.create(client, request);
    }

    @PreAuthorize("hasRole('CLIENT')")
    @PostMapping(path = "/update")
    public ResponseEntity<?> updateClient(@RequestBody Client client, HttpServletRequest request) {
        log.info("Ip: {}, Username: {}, Client is successfully updated!",request.getRemoteAddr(), client.getUsername());
        return clientService.updateClient(client , request);
    }

    @GetMapping(path = "/activate")
    public ResponseEntity<?> activateClientAccount(HttpServletRequest request, @RequestParam("token") String hashCode) {
        ConfirmationToken token = confirmationTokenRepository.findByConfirmationToken(hashCode);
        Long secs = (token.getCreatedDate().getTime() - new Date().getTime())/1000;
        Client verificationClient = token.getClient();
        if (verificationClient == null || verificationClient.isActivated() || secs > 3600 ) {
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        }

        verificationClient.setActivated(true);
        log.info("Ip: {}, Username: {}, Client is successfully activated!",request.getRemoteAddr(), verificationClient.getUsername());
        clientService.save(verificationClient);

        HttpHeaders headers = new HttpHeaders();
        headers.add("Location", "https://localhost:4200");
        return new ResponseEntity<String>(headers, HttpStatus.OK);
    }

    @PutMapping(path = "/newPassword/{username}")
    public ResponseEntity<?> sendNewPassword(@PathVariable String username) {
        if(clientService.findByUsername(username) != null)
            return clientService.sendNewPassword(clientService.findByUsername(username));
        if(companyService.findByUsername(username) != null)
            return companyService.sendNewPassword(companyService.findByUsername(username));
        return new ResponseEntity<String>(HttpStatus.BAD_REQUEST);
    }

    @PreAuthorize("hasRole('CLIENT')")
    @GetMapping(path = "/username/{username}")
    public ResponseEntity<?> clientByUsername(@PathVariable String username) {
        return clientService.getByUsername(username);
    }
}
