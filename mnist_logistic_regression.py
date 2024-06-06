import torch
import torch.nn as nn
from torchvision import datasets, transforms

class Logistic_Regression(nn.Module):
	def __init__(self):
		super().__init__()
		self.linear = nn.Linear(28*28, 10)

	def forward(self, x):
		return self.linear(x.view(-1, 28*28))

	def evaluate(self, test_dataloader):
		with torch.no_grad():
			total = 0
			correct = 0 
			for data in test_dataloader:
				inputs, targets = data
				outputs = self.forward(inputs)
				_, predicted = torch.max(outputs.data, 1)
				total += targets.size(0)
				correct += (predicted == targets).sum().item()
		return 100*(correct/total)

batch_size = 100
epochs = 10
lr = 0.1

model = Logistic_Regression()
criterion = nn.CrossEntropyLoss()
opt = torch.optim.SGD(model.parameters(), lr = lr)


training_dataset = datasets.MNIST(root='./data', train=True, transform=transforms.ToTensor(), download=True)
test_dataset = datasets.MNIST(root='./data', train=False, transform=transforms.ToTensor(), download=True)

training_dataloader = torch.utils.data.DataLoader(training_dataset, batch_size=batch_size, shuffle=True)
test_dataloader = torch.utils.data.DataLoader(test_dataset, batch_size=batch_size, shuffle=False)

for epoch in range(epochs):
	print('Epoch', epoch, '- Test Accuracy', model.evaluate(test_dataloader))
	for data in training_dataloader:
		opt.zero_grad()
		inputs, targets = data
		outputs = model(inputs)
		loss = criterion(outputs, targets)
		loss.backward()
		opt.step()
	
print("Final Accuracy", model.evaluate(test_dataloader))


torch.save(model.state_dict(), "./mnist_logistic_regression.pt")