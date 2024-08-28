import React, { useState } from 'react';
import Loader from './Loader';
import { Box, InputGroup, IconButton, InputRightElement, Button, FormControl, FormLabel, Input, Heading, Text } from '@chakra-ui/react';
import { BASE_URL } from '../App';
import { FaEye, FaEyeSlash } from "react-icons/fa";

const Login: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    
    try {
      const response = await fetch(BASE_URL + 'login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        throw new Error('Email ou senha inválidos');
      }

      const data = await response.json();

      localStorage.setItem('token', data.token);

      window.location.href = '/';
    } catch (err:any) {
      setError(err.message);
      console.error('Erro ao fazer login:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleClick = () => setShowPassword(!showPassword);

  return (
    <Box
      p={4} 
      borderWidth={1} 
      borderRadius="lg" 
      display="flex" 
      flexDirection="column" 
      justifyContent="center" 
      alignItems="center"
    >
      {loading && <Loader />}
      <Heading as="h2" size="lg" mb={4}>Login</Heading>
      <form onSubmit={handleSubmit}>
        <FormControl mb={4} isRequired>
          <FormLabel>Email</FormLabel>
          <Input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Digite seu email"
            w={300}
            required
          />
        </FormControl>
        <FormControl mb={4} isRequired>
          <FormLabel>Senha</FormLabel>
          <InputGroup>
            <Input
              type={showPassword ? 'text' : 'password'}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Digite sua senha"
              w={300}
              required
            >
            </Input>
            <InputRightElement>
              <IconButton
                variant="link"
                aria-label={showPassword ? 'Ocultar senha' : 'Mostrar senha'}
                icon={showPassword ? <FaEye /> : <FaEyeSlash />}
                onClick={handleClick}
              />
            </InputRightElement>
          </InputGroup>
        </FormControl>
        <Button type="submit" colorScheme="teal" width="full">Entrar</Button>
      </form>
      <Text my={5}>Ainda não tem uma conta? <a href="/register" style={{color: 'teal'}}>Registre aqui</a></Text>
      {error && <Text my={5} color="red.500">{error}</Text>}
    </Box>
  );
};
  
export default Login;