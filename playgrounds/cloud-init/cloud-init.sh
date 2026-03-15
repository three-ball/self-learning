#!/bin/bash

VM_USER="vienct"
VM_PASSWORD="Banhmi123!!" # May use hash?
VM_USER_PUBKEY="ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC8FvF9/CUYwG9jgaFc7ZWmhdTw2V7l5m//jTQPXaTHB6d1nv4PgbzppFd6+rb9toR1mKggNv+d4aHNqpgSxcMAufcN5Wi1yKpeo6/TKvBJJG3gM2LFlSIAUztXjMxxWIe2Pbn+UBudxjVNMXsqZAQ/Jnl1DnXP+BO7GDFNNUbhgJjM2vZZOjx62XOPnX7dPX28yDSNrY25J4xDoU2uxkf/VjEkAd107iDA2gFyqyf0GGX2f0SXeK0dX9Y7jMWG3IRziVLw3s4jUxP6YyCXV57Qrk5bD09gPQuHeBLGrwS9wxltIEkRhk9wog6ww8LgIw5TkNsBFnfWcYUH2nrOrkJm7g8oxGWOcPhe4LYdVBFpeMUcN4+Q4oBfO/l5LBXDliE8/C9+jyWXkiTbCUiu6lJ8WxotPPLwXhkS9w3NUL57uXwfAWbgSkQrNcSZ1LeWkeoCCd81oaB/F60qoLGE28s8pCGcBICF/o16VYHcpoh8r938yaJtunX357/N267hmy8= vienct@DESKTOP-LJLQ3N1"
VM_NAME="vienct3-rocky10-prod"
VM_IMAGE="${VM_NAME}.qcow2"
VM_USERDATA="${VM_NAME}-userdata.yml"
VM_NETWORK="${VM_NAME}-network.cfg"
VM_STATIC_IP="192.168.250.185/24"

# Download the latest Rocky Linux 10 cloud image
wget https://dl.rockylinux.org/pub/rocky/10/images/x86_64/Rocky-10-GenericCloud-Base.latest.x86_64.qcow2

# Make a copy of the downloaded image for vienct3
cp --reflink Rocky-10-GenericCloud-Base.latest.x86_64.qcow2 ${VM_IMAGE}

# or using qemu-img
qemu-img create -f qcow2 -F qcow2 -b Rocky-10-GenericCloud-Base.latest.x86_64.qcow2 ${VM_IMAGE}

# Resize the copied image to 50GB
qemu-img resize ${VM_IMAGE} 50G

# Meta-Data tells cloud-init where it is. User-Data tells cloud-init what to do.
# When troubleshooting, always check /var/log/cloud-init.log

# Create a cloud-init configuration file
# Create a simple user-data.yml
cat <<EOF > ${VM_USERDATA}
#cloud-config
hostname: ${VM_NAME}
users:
  - name: ${VM_USER}
    groups: users, admin
    sudo: ALL=(ALL) NOPASSWD:ALL
    lock_passwd: false
    plain_text_passwd: ${VM_PASSWORD}
    ssh_authorized_keys:
      - ${VM_USER_PUBKEY}
EOF

cat <<EOF > ${VM_NETWORK}
#cloud-config
network:
  version: 2
  ethernets:
    enp1s0:
      dhcp4: no
      addresses: [${VM_STATIC_IP}]
      gateway4: 192.168.250.1
      nameservers:
        addresses: [8.8.8.8]
EOF


virt-install --name ${VM_NAME} \
--ram 4096 \
--vcpus 4 \
--disk path=${VM_IMAGE},format=qcow2 \
--cloud-init user-data=/root/libvirt-images/${VM_USERDATA},network-config=/root/libvirt-images/${VM_NETWORK} \
--network bridge=br-qemu,model=virtio \
--os-variant rocky10 \
--import --noautoconsole



virsh console ${VM_NAME}