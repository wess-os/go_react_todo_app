import React, { useState } from 'react';
import Loader from './Loader';
import { Box, InputGroup, IconButton, InputRightElement, Button, FormControl, FormLabel, Input, Heading, Text } from '@chakra-ui/react';
import { BASE_URL } from '../App';
import { FaEye, FaEyeSlash } from "react-icons/fa";

const Register: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [showPassword, setShowPassword] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      const response = await fetch(BASE_URL + 'register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, email, password }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(await data.error);
      }

      setSuccess('Registro bem-sucedido!');

      setTimeout(() => {
        window.location.href = '/login';
      }, 2000);
    } catch (err: any) {
        setError(err.message);
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
      <Heading as="h2" size="lg" mb={4}>Registro</Heading>
      <form onSubmit={handleSubmit}>
        <FormControl mb={4} isRequired>
          <FormLabel>Email</FormLabel>
          <Input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Digite seu email"
            w={300}
          />
        </FormControl>
        <FormControl mb={4} isRequired>
          <FormLabel>Usuário</FormLabel>
          <Input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Informe seu nome"
            w={300}
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
            />
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
        <Button type="submit" colorScheme="teal" width="full">Registrar</Button>
      </form>
      <Text my={5}>Já tem uma conta? <a href="/login" style={{color: 'teal'}}>Entre aqui</a></Text>
      {error && <Text my={5} color="red.500">{error}</Text>}
      {success && <Text my={5} color="green.500">{success}</Text>}
    </Box>
  );
};

export default Register;